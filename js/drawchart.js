  function updateChartData() {
    setInterval(function() {
      var url = "http://localhost:3000/data/";
        xmlhttp=new XMLHttpRequest();
       
        xmlhttp.onreadystatechange=function()
        {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)
          {
            postMessage(xmlhttp.responseText);
            
          }
        }
        xmlhttp.open("GET", url, false);
        xmlhttp.send();

      }, 1000);
  }
  updateChartData();