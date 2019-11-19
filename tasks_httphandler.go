package task

import (
	"fmt"
	"net/http"
	"strings"
)

// CounterOutputHttpHandler Http Handler for output counter info
func (service *TaskService) CounterOutputHttpHandler(w http.ResponseWriter, r *http.Request) {
	tableData := ""
	for _, v := range service.GetAllTaskCountInfo() {
		tableData += "<tr><td>" + v.TaskID + "</td><td>" + v.Lable + "</td><td>" + fmt.Sprint(v.Count) + "</td></tr>"
	}
	col := `<colgroup>
		  <col width="40%">
		  <col width="30%">
          <col width="30%">
		</colgroup>`
	header := `<tr>
          <th>TaskID</th>
      	  <th>Lable</th>
          <th>Value</th>
        </tr>`
	html := createOneTableHtml(col, "TaskCounterInfo", header, tableData)
	w.Write([]byte(html))
}

// TaskOutputHttpHandler Http Handler for output task info
func (service *TaskService) TaskOutputHttpHandler(w http.ResponseWriter, r *http.Request) {
	tableData := ""
	for _, v := range service.GetAllTasks() {
		tableData += "<tr><td>" + v.TaskID() + "</td><td>" + v.GetConfig().TaskType + "</td><td>" + fmt.Sprint(v.GetConfig().IsRun) + "</td><td>" + v.GetConfig().Express + "</td><td>" + fmt.Sprint(v.GetConfig().DueTime) + "</td><td>" + fmt.Sprint(v.GetConfig().Interval) + "</td></tr>"
	}
	col := `<colgroup>
		  <col width="20%">
		  <col width="10%">
          <col width="10%">
		  <col width="30%">
		  <col width="15%">
		  <col width="15%">
		</colgroup>`
	header := `<tr>
          <th>TaskID</th>
      	  <th>TaskType</th>
          <th>IsRun</th>
          <th>Express</th>
          <th>DueTime</th>
          <th>Interval</th>
        </tr>`
	html := createOneTableHtml(col, "TaskInfo", header, tableData)
	w.Write([]byte(html))
}

func createOneTableHtml(col, title, header, body string) string {
	template := `<br><table class="bordered">
		{{col}}
		<caption>{{title}}</caption>
	  <thead>
	 {{header}}
	  </thead>
		{{body}}
	</table>`
	data := strings.Replace(template, "{{col}}", col, -1)
	data = strings.Replace(data, "{{title}}", title, -1)
	data = strings.Replace(data, "{{header}}", header, -1)
	data = strings.Replace(data, "{{body}}", body, -1)
	html := strings.Replace(tableHtml, "{{tableBody}}", data, -1)

	return html
}

var tableHtml = `<html>
<html><head>
   <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0;">
 
  <meta name="Generator" content="EditPlus">
  <meta name="Author" content="">
  <meta name="Keywords" content="">
  <meta name="Description" content="">
  <title>DotTask</title>
    <style>
    .overtable {
      width: 100%;
      overflow: hidden;
      overflow-x: auto;
    }
    body {
      max-width: 780px;
       margin:0 auto;
      font-family: 'trebuchet MS', 'Lucida sans', Arial;
      font-size: 1rem;
      color: #444;
    }
    table {
      font-family: 'trebuchet MS', 'Lucida sans', Arial;
      *border-collapse: collapse;
      /* IE7 and lower */
      border-spacing: 0;
      width: 100%;
      border-collapse: collapse;
      overflow-x: auto
    }
    caption {
      font-family: 'Microsoft Yahei', 'trebuchet MS', 'Lucida sans', Arial;
    text-align: left;
    padding: .5rem;
    font-weight: bold;
    font-size: 110%;
    color: #666;
    }
    tr {
      border-top: 1px solid #dfe2e5
    }
    tr:nth-child(2n) {
      background-color: #f6f8fa
    }
    td,
    th {
      border: 1px solid #dfe2e5;
    }
    .bordered tr:hover {
      background: #fbf8e9;
    }
    .bordered td,
    .bordered th {      padding: .6em 1em;

      border: 1px solid #ccc;
      padding: 10px;
      text-align: left;
    }
  </style>
  <script>
  
(function(doc, win) {
    window.MPIXEL_RATIO = (function () {
        var Mctx = document.createElement("canvas").getContext("2d"),
            Mdpr = window.devicePixelRatio || 1,
            Mbsr = Mctx.webkitBackingStorePixelRatio ||
                Mctx.mozBackingStorePixelRatio ||
                Mctx.msBackingStorePixelRatio ||
                Mctx.oBackingStorePixelRatio ||
                Mctx.backingStorePixelRatio || 1;
    
        return Mdpr/Mbsr;
    })();

    function addEventListeners(ele,type,callback){
    
        try{  // Chrome、FireFox、Opera、Safari、IE9.0及其以上版本
            ele.addEventListener(type,callback,false);
        }catch(e){
            try{  // IE8.0及其以下版本
                ele.attachEvent('on' + type,callback);
            }catch(e){  // 早期浏览器
                ele['on' + type] = callback;
            }
        }
    }

    var docEl = doc.documentElement,
        resizeEvt = 'orientationchange' in window ? 'orientationchange' : 'resize';
    window.recalc = function() {
            var clientWidth = docEl.clientWidth < 768 ? docEl.clientWidth : 768;
            if (!clientWidth) return;
            docEl.style.fontSize = 10 * (clientWidth / 320) *  window.MPIXEL_RATIO + 'px';
        };
    window.recalc();
    
    addEventListeners(win, resizeEvt, recalc);
})(document, window);

</script>
</head>
<body>
<div class="overtable">
{{tableBody}}
</div>
</body>
</html>
`
