var service = 300;
// Load the Visualization API library and the piechart library.
// google.load('visualization', '1.0', {'packages':['corechart']});
// google.setOnLoadCallback(drawChart);

function setupGoogleCharts() {
  google.load('visualization', '1.0', {'packages':['corechart']});
  google.setOnLoadCallback(drawChart);
}

function drawChartT2T3() {
  data = new google.visualization.DataTable();
  data.addColumn('number', 'Packet Number');
  data.addColumn('number', 'Latency');
  data.addColumn('number', 'Service');

  var options = {
    title: '',
    hAxis: {
      title: 'Packet Number',
      slantedText: false,
      viewWindow: {
        max: data.getColumnRange(0)["max"],
        min: data.getColumnRange(0)["min"]
      },
    },
    vAxis: {
      title: 'Latency (milliseconds)',
      slantedText: false,
      max: 10,
      viewWindow: {
        min: 0
      }
    },
    legend: 'none',
    interpolateNulls: true,
    series: {
      1: { lineWidth: 1, pointSize: 0 }
    },
    sortAscending: true
  };

  var chart = new google.visualization.ScatterChart(document.getElementById('chartSecond_div'));
  chart.draw(data, options);
  console.log(data.getColumnRange(0)["max"]);
}

function addDataToChartT2T3(number, latency) {
  data.addRows([[number, latency]]);
  var options = {
    title: '',
    hAxis: {
      title: 'Packet Sequence Number',
      viewWindow: {
        max: data.getColumnRange(0)["max"],
        min: data.getColumnRange(0)["min"]
      },
    },
    vAxis: {
      title: 'Latency (ms)',
      max: 10,
      viewWindow: {
        min: 0
      }
    },
    legend: 'none',
    interpolateNulls: true,
    series: {
      1: { lineWidth: 1, pointSize: 0 }
    },
    sortAscending: true
  };

  var chart = new google.visualization.ScatterChart(document.getElementById('chartSecond_div'));
  chart.draw(data, options);
}

function setServiceLevelOnChart(serviceLevel) {
  // service = serviceLevel
}

export {addDataToChart, drawChart, setServiceLevelOnChart}
