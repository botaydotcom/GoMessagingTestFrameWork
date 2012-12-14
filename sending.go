func performTestCaseOnce(addr *net.TCPAddr, testCase *TestCase, resultChan chan TestResult) {
	var result = TestResult{
		IsCorrect: true,
		Reason:    nil}

	var listConnection = make(map[int]*net.TCPConn)
	var connectionState = make(map[int]bool)
	for i := range connectionState {
		connectionState[i] = false // all closed
	}

	var totalTime int64 = 0
	for i, message := range testCase.MessageSequence {
		// parse message (plug in values if necessary)
		parsedMessage, comm, useBase, baseCmd, err := parseAMessage(message)
		if err != nil {
			result.IsCorrect = false
			errMsg := fmt.Sprintf("Cannot parse message %d", i)
			result.Reason = errors.New(errMsg)
			break
		}
		protoParsedMessage := parsedMessage.(proto.Message)

		// get connection from list or make one if it does not exist yet
		connectionId := message.Connection
		if isOpened, exist := connectionState[connectionId]; (!exist ||  !isOpened) {
			//Connection does not exist or closed, create new connection for connection 
			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				result.IsCorrect = false
				result.Reason = err
				break
			} else {
				listConnection[connectionId] = conn
				connectionState[connectionId] = true
			}
		} else {
			// Connection  exists => reuses
		}
		conn := listConnection[connectionId]

		if message.FromClient { // message from client, so send to server now
			begin := time.Now()
			sendMsg(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)
			end := time.Now()
			totalTime += end.Sub(begin).Nanoseconds()
		} else { // message from server, so read from buffer now
			
			if (DEBUG_READING_MESSAGE) {
				fmt.Println("Expecting to receive a message from server ", message.MessageType)
			}
			
			replyMessage, err := readReply(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)
			if err == nil {
				if comparedResult, err := compareGetValueForProtoMessage(protoParsedMessage, *replyMessage); comparedResult {
					// CORRECT Reply message.
				} else {
					// INCORRECT Reply message 
					result.IsCorrect = false
					errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
					result.Reason = errors.New(errMsg)
					break
				}
			} else {
				// Encounter error while trying to parse reply message
				result.IsCorrect = false
				errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
				result.Reason = errors.New(errMsg)
				break
			}
		}
		if message.Close {
			connectionState[connectionId] = false
			conn.SetLinger(DEFAULT_LINGER_PERIOD)

			if (conn.Close() == nil) {
				fmt.Println("CONNECTION SUCCESSFULLY CLOSED")
			}
		}
	}

	for i := range listConnection {
		if connectionState[i] {
			listConnection[i].Close()
		}
	}

	result.TimeTaken = totalTime

	cleanUpAfterTest(addr, testCase)

	resultChan <- result
}

/***************Send and Receive Messages********************************/
// send a message with command and message using given connection
func sendMsg(useBase bool, baseCmd byte, cmd byte, msg proto.Message, conn *net.TCPConn) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	length := int32(len(data)) + 1

	if (useBase) {
		length = length + 1
	}

	if (DEBUG_SENDING_MESSAGE) {
		log.Print("sending message base length: ", length, " command / command / message: ", baseCmd, cmd, msg)
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, length)
	if (useBase) {
		binary.Write(buf, binary.LittleEndian, baseCmd)
	}

	binary.Write(buf, binary.LittleEndian, cmd)
	buf.Write(data)
	// write data to connection
	numSent, err := conn.Write(buf.Bytes())
	if (err != nil) {
		log.Fatal("error while sending data to server ", err, " num bytes sent: ", numSent)
	} else {
		fmt.Println("MESSAGE SUCCESSFULLY SENT", numSent)
	}
}



// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(useBase bool, expBaseCmd byte, expCmd byte, expMsg proto.Message, conn *net.TCPConn) (*proto.Message, error) {
	length := int32(0)
	duration, err := time.ParseDuration(READ_TIMEOUT)
	timeNow := time.Now()
	// set timeout
	err = conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		return nil, err
	}
	// read from buffer
	err = binary.Read(conn, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	// check if use base-command
	if (useBase) {
		length = length - 1
		var baseCmd byte
		err = binary.Read(conn, binary.LittleEndian, &baseCmd)

		if (baseCmd != expBaseCmd) {
			// finish reading the rest of the message, 
			// so that it does not affects other message
			rbuf := make([]byte, length)
			io.ReadFull(conn, rbuf)
			// report error
			errMsg := fmt.Sprintf("Unexpected BASE CMD %d", baseCmd)
			return nil, errors.New(errMsg)	
		}
	}

	rbuf := make([]byte, length)
	_, err = io.ReadFull(conn, rbuf)
	if err != nil {
		return nil, err
	}

	var res proto.Message
	if len(rbuf) < 1 {
		errMes := fmt.Sprintf("Reply message is too short: %d", len(rbuf))
		return nil, errors.New(errMes)
	}

	cmd := rbuf[0]
	// command needs to be equal to expected command
	if cmd != expCmd {
		errMsg := fmt.Sprintf("Unexpected CMD %d", cmd)
		return nil, errors.New(errMsg)
	}
	
	newValue := reflect.New(reflect.ValueOf(expMsg).Elem().Type()).Interface()
	res = newValue.(proto.Message)
	err = proto.Unmarshal(rbuf[1:], res)

	if err != nil {
		log.Fatal(err)
	}

	if (DEBUG_READING_MESSAGE) {
		log.Print("Successfully receive message from network: ", res)
	}
	return &res, err
}
