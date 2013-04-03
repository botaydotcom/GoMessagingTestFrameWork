GoMessagingTestFrameWork
========================

A Test framework for Btalk test server


Running Test:
========================
To run test:

	"autoTest.exe (-inputDir inputDir) (-timeOut timeOutForReadOperation)"

Test suite description:
========================

Structure of a test suite:
---------------------

	<TestSuite>
		<TestSuiteName>Chat</TestSuiteName>
		<TargetHost>203.117.155.188</TargetHost>
		<TargetPort>9100</TargetPort>
		<VarMap>(repeated <Var></Var>)</VarMap> (optional)
		<IgnoreMessages>(repeated <Message></Message>)</IgnoreMessages> (optional, global ignore message list)
		<ListTest>
			(repeated <TestInfo></TestInfo>)
		</ListTest>
		<Tests>
			(repeated <Test></Test>)
		</Tests>
	</TestSuite>

The test suite consists of the following sections:

0. Test suite info:
---------------------
	<TestSuiteName>Chat</TestSuiteName>
	<TargetHost>203.117.155.188</TargetHost>
	<TargetPort>9100</TargetPort>


1. Global var map:
---------------------
#####- Purpose:
In the test suite, sometimes some values can be repeated a lot of times. We want to introduce variables to help here:
#####- How to add:
To specify the list of varible to be used in the whole test suite, and specify their values, add section <VarMap></VarMap>
#####For example:

	<VarMap>
		<!-- to define a var name: email1, with value: indotest15@test.com-->
		<Var name="email1">
			<Value>indotest15@test.com</Value>
		</Var>	
		(.. repeat <Var> </Var> for other variables / values)
	</VarMap>

Inside the test suite, to use a variable, use the notion: {{.varname}}. If varname is a variable that's not defined in varmap, it value will be undertood as: "waiting to be filled". The first time this variable getting assigned (by value returned from server), it will carry that value from then on.

This can be used in case such as: Getting user id from server (assign to an unbound variable), then later on use that variable to use that userid value.

2. Ignore list:
---------------------
#####- Purpose:
For some test cases, some type of messages can be missing / arrive out of order (especially for buddy online messages).

Sometimes, this is "expected", and some unimportant messages, such as buddy_online can be ignored when that problem happens. We can add them to "ignore list"
#####- How to add:
Adding the section 

	<IgnoreMessage></IgnoreMessage>:

#####For example:

	<IgnoreMessages>
		<!-- to ignore buddy online message (base command = 0x00, command = 115 = 0x73)-->
		<Message>
			<BaseCommand>0x00</BaseCommand>
			<Command>0x73</Command>
			<MessageType>Auth_Buddy_S2C_BuddyOnline</MessageType>
		</Message>
		(... repeat <Message></Message> for other types of ignored messages)
	</IgnoreMessages>

This is the global ignore messages list (to be ignored for any test in the whole test suite). This section can also be added to each of the <TestInfo> tag to specify additional ignore message types for the specific test:

#####For example:

	<ListTest>
		<TestInfo skip="false" repeat="1">
			<Name>A and B log in, A sends a message to B (both online), B receives</Name>
			<IgnoreMessages>
				<Message>
					<BaseCommand>0x00</BaseCommand>
					<Command>0x73</Command>
					<MessageType>Auth_Buddy_S2C_BuddyOnline</MessageType>
				</Message>
				(... repeat <Message></Message> for other types of ignored messages)
			</IgnoreMessages>
		</TestInfo>
		(..repeated <TestInfo></TestInfo>)
	</ListTest>

3. Test info section:
---------------------
#####- Purpose:
To specify information specify to each test case (name/ how many times to run / skipped or not...)
#####- How to add:
Add the tag: 

	<ListTest>
		(repeated <TestInfo></TestInfo> for each of the test case)
	</ListTest>

For each <TestInfo> tag:

	<TestInfo skip="false" repeat="1">
		<Name>A and B log in, A sends a message to B (both online), B receives</Name>
		(optional <IgnoreMessages></IgnoreMessages> as specified in the previous section)
	</TestInfo>

######skip : whether to skip the test case.
######repeat : how many times to repeat the test case. Currently, all of the iterations will be run CONCURRENTLY.

	<Name></Name>
######name of the test case

#####For example:

	<ListTest>
		<!--specify a test case that will not be skipped, is run only once, with a specific name-->
		<TestInfo skip="false" repeat="1">
			<Name>A and B log in, A sends a message to B (both online), B receives</Name>
		</TestInfo>
	</ListTest>
	
4. Test messages section:
---------------------
#####- Purpose:
To specify the lists of messages to be sent / received for each test case.
#####- How to add:
Add the tag: 

	<Tests>
		(repeated <Test></Tests> for each of the test case)
	</Tests>

Each <Test> tag specifies the message sequence and clean up message sequence for a test case. Repeat for each test case in the suite.

Inside <Test> tag:

	<Test>
		<MessageSequence>
			(repeated <Message></Message>)
		</MessageSequence>
		<CleanUpSequence>
			(repeated <Message></Message>)
		</CleanUpSequence>
	</Test>

The <MessageSequence> tag specifies the message sequence for a test case. 

Inside <MessageSequence> tag:

	<MessageSequence>
		(repeated <Message></Message> for each message of the sequence)
	</MessageSequence>

Inside the <MessageSequence>, the messages are specified in order. All the messages will be sent to / received from server one after another in order.

Each <Message> is specified as:

	<Message type="Auth_C2S_LoginInfo" fromClient="true" connection="1" close="true">		
		<BaseCommand>(basecommand)</BaseCommand> (optional)
		<Command>(command)</Command>
		<ForceCheck>
			<FieldName>(field name specified in domain format)</FieldName>
		</ForceCheck>(optional)
		<Data>
			(raw data in xml format of the message)
		</Data>
	</Message>

######type: type of the messages, in format: (packagename_typename)
######fromClient: true/false
######connection: multiple connections can be opened during a test (for example: A and B both logs in). This field specify the connections that the message belongs to.
######close: true/false. true to force close this connection (logout) after finish processing this message.

	<BaseCommand>(base command of message)</BaseCommand> (optional)
	<Command>(command of the message)</Command>
	<Data>
		(raw data of the message in xml format)
	</Data>

######ForceCheck: normally, for a message from server, the test framework only checks for value of the field that exists in the specification of the expected message, i.e: if the field is not specified in the expected message, the field is ignored. For instance: if a message can optionally have a field named: gender, but in the expected message specification, we leave out this field, this field is not checked in the message returned from server. This is used to ignore trivia fields that are not important.
However, sometime, this behaviour is not desired, especially when there's an array that we need to make sure it's empty. To force the test framework to check for a field even though it's not specified (= need to be empty), we can include the flag: ForceCheck in format:
	
	<ForceCheck>
		<FieldName>(field name specified in domain format)</FieldName> (repeated)
	</ForceCheck>

Repeated <FieldName> specifies the name of the fields that need to be checked. 

Domain format: is a special format to specify a field inside a message. The field name always starts with a ".". After that, the name of the fields (include nested field) that lead to the specified field is concatenated, separated by "."

For example: In a message with structure:

	<A>
		<B>
			<C>
			</C>
		</B>
	</A>

to specify field C, the domain name for C is:
".B.C" (A is the message name, which can be ignored). The domain name for B is ".B".

######ByteEncoded:
In some message, a field could be the value of a wholly byte-encoded message. For example: In the binding account message, the field OAuthRawInfo is encoded in this manner. 
In this case, we cannot use a helper function to byte-encode the message and inject the value to the field. The reason is because after being byte-encoded, the resulted data may contain special character that are not understood by the XML parser. In this case, we need to use ByteEncoded field:

	<ByteEncoded>
		<Var type="Auth_C2S_OAuthRawInfo" >
			<Field>.OAuthInfo</Field>
			<Data>
				(a message description here)
			</Data>
		</Var> (repeated Var for multiple byte-encoded message)
	</ByteEncoded>

Each byte-encoded message specification starts with tags Var. The Type field refers to the raw message type that will be byte-encoded. The Field tag contains a value written in domain-name format (mentioned above) that specifies the field in the following message that this byte-encoded message will be injected into. The Data tag contains the specification for the message itself. Format of the message is just like any other message.


#####Note: 
Inside the raw data section, we can use {{.varname}} notion to use variable (as specified in the var map section)

#####For example:

	<Message type="Auth_S2C_LoginUserInfo" fromClient="false">
		<Command>0x02</Command>
		<Data>
			<Auth_S2C_LoginInfoResult>
				<MyInfo>
					<UserId>{{.userId1}}</UserId>
					<Name>{{.email1}}</Name>						
					<NickName>{{.name1}}</NickName>
				</MyInfo>
			</Auth_S2C_LoginInfoResult>
		</Data>
	</Message>

The <CleanUpSequence> tag specifies the clean up message sequence for a test case. This is a sequence of messages to open up (logins) some connections for some users, and let the program perform message clean up (ack all pending offline message) for that connection.

This part is retained for historical reason, but its function is better replace by using cleanUpEngine.

Clean up:
========================
Clean up failed test:
The autotest program will run non-stop, repeating all test cases in the testing directory. If it encounters any error in any test case, it will stop at the test case.

In that case, some connections will be left open, and some states on server can be changed. To clean up this so that we can re-run the test freshly, run: 

	"cleanUpEngine.exe (-in cleanUpFileName)"

The cleanUpFile has the CleanUpSequence, with the same meaning as the CleanUpSequence mentioned bove.

List of files in this directory:
========================
#####Folders: 
All folders with single .pb.go files are the protobuf messages (in Go)
Opcode folder: file that specifies operation code for the messages
proto files: original proto specification files.

js: folder that contains javascript and css file for the displayResult.html page (display a monitoring page for stresstest)
IndividualCases: individual test specification, can be used for feature test.
Mics go files: some miscellaneous source files.
Documentation: documentation.

#####Feature test: 
xmltest.go : standard feature test (that runs single test case). Default test file is test.xml. For more option refer to documentation.
cleanUp.go : standard clean up program. Refer to documentation for more information.
generateMapCode.go : generate a map of protobuf messages that can be used to create protobuf message dynamically.
generateMapFuncCode.go : generate a map of helper functions that can used in the test specification. The helper function are defined in the file helper.func.go in folder Helper.

#####Auto test:
autoTest.go : autotest program (runs multiple test cases, continuously). By default, autotest program will run all cases that are stored in folder TestCase.
Refer to documentation for more information.
TestCase: folder that contains xml file (test specifications) that can be used for autoTest engine.
TestEngine: autoTest engine.

#####Stress test:
stressTest.go : standard stresstest program. (spawning multiple go routines to stress test server using a set of specifications).

StressTestCase: folder that contains a set of test specifications that are used in the stressTest engine.
StressTestEngine* : different stress test engine. Refer to stressTest.go for more information



