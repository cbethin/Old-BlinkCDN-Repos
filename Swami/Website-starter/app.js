// import {MDCButton} from '@material/button';
import {MDCMenu} from '@material/menu';
import {MDCChipSet} from '@material/chips';
import axios from 'axios';

import {makeCORSRequest} from './src/CORS';
import {getRealTimeUpdatesTotal, getRealTimeUpdatesFirst, getRealTimeUpdatesSecond} from './src/firebaseCode';
import {addDataToChart, addDataToChartT1T2, addDataToChartT2T3, setServiceLevelOnChart, updateRegionsMap} from './src/charts';

// Constants
// const SERVER_ADDR = "155.246.45.43:8080"
const SERVER_ADDR = "18.221.128.102:8080"
// const CLIENT_ADDR = "155.246.113.15:8000"
const CLIENT_ADDR = "18.216.192.154:8081"

const addresses = ["52.53.177.194:8001", "35.176.239.10:8001", "18.184.225.196:8001", "13.115.224.27:8001"];
// const addresses = ["155.246.45.43:8001", "155.246.45.43:8002", "155.246.45.43:8003", "155.246.45.43:8004"]
const locations = ["United States", "Germany", "United Kingdom", "Japan"];
var isMapShowing = false;
var srcLocation, dstLocation, selectedPath;


$(document).ready(() => {

  sendStopTest();
  window.addEventListener("beforeunload", function(e) {
    sendStopTest();
  });

  // Elements
  const sourceChipSet = new MDCChipSet(document.querySelector('#src-chipset'));
  const dstChipSet = new MDCChipSet(document.querySelector('#dest-chipset'));
  const nextBtn = document.querySelector('#next-btn');
  const toGraphBtn = document.querySelector('#to-graph-btn');
  const showMapsBtn = document.querySelector('#show-map-btn');

  var selectedIndex = 0;

  nextBtn.addEventListener('click', () => {
    console.log("CLICK");
    var srcAddr, dstAddr;
    for (var chip in sourceChipSet.chips) {
      if (sourceChipSet.chips[chip].isSelected()) {
        srcAddr = addresses[chip];
        srcLocation = location[chip];
      }
    }
    for (var chip in dstChipSet.chips) {
      if (dstChipSet.chips[chip].isSelected()) {
        dstAddr = addresses[chip];
        dstLocation = location[chip];
      }
    }

    makeCORSRequest('GET', `http://${SERVER_ADDR}/getSLA?source=${srcAddr}&destination=${dstAddr}`, (response) => {
      transitionToSelectSla(response);
    });
  });

  toGraphBtn.addEventListener('click', () => {
    setServiceLevelOnChart(parseFloat($('.latency-label')[selectedIndex].innerHTML.split(' ')[0]))
    makeCORSRequest('POST', `http://${SERVER_ADDR}/setSLA?sid=0001000100010001&path=${selectedIndex}`, transitionToGraphs);
  });

  showMapsBtn.addEventListener('click', () => {
    if (!isMapShowing) {
      $('#graph-container1').css('height', '0px');
      $('#graph-container1').css('opacity', '0');
      $('#graph-container2').css('height', '500px');
      $('#graph-container2').css('opacity', '1');
      $('#graph-label').html("Activation of Nodes in Path")
    } else {
      $('#graph-container1').css('height', '700px');
      $('#graph-container1').css('opacity', '1');
      $('#graph-container2').css('height', '0px');
      $('#graph-container2').css('opacity', '0');
      $('#graph-label').html("End-To-End Path Diagnostics")
    }

    isMapShowing = !isMapShowing;
  })

  $('.card-container__card').click(() => {
    var cards = $('.card-container__card');
    for (var i = 0; i < cards.length; i++) {
      if (cards[i] == $(event.target)[0]) {
        $(cards[i]).addClass('selected');
        selectedIndex = i;
      } else {
        $(cards[i]).removeClass('selected');
      }
    }
  })

  $('#stop-test-btn').click(() => {
    $('#stop-test-btn').prop('disabled', true);
    sendStopTest()
  })

  $('#brand').click(() => {
    window.location.reload();
  })
});

function transitionToGraphs(response) {
  var pathLocations = response.split(' ');

  var usActivation = 0;
  var gmActivation = 0;
  var ukActivation = 0;
  var jpActivation = 0;

  if (pathLocations.includes(addresses[0])) {
    usActivation = 1;
  }
  if (pathLocations.includes(addresses[1])) {
    gmActivation = 1;
  }
  if (pathLocations.includes(addresses[2])) {
    ukActivation = 1;
  }
  if (pathLocations.includes(addresses[3])) {
    jpActivation = 1;
  }

  updateRegionsMap(usActivation, gmActivation, ukActivation, jpActivation);

  if (response != undefined) {
    $('.SLA-page').animate({
      left: "-100vw",
      opacity: 0
    }, 300);
    $('.graph-page').animate({
      left: "0",
      opacity: 1
    }, 300);
  }

  getRealTimeUpdatesTotal((data) => {
    // console.log(data);
    addDataToChart(data.packetNumber, data.time*1000.0)
  })

  getRealTimeUpdatesFirst((data) => {
    // console.log(data);
    addDataToChartT1T2(data.packetNumber, data.time*1000.0)
  })

  getRealTimeUpdatesSecond((data) => {
    // console.log(data);
    addDataToChartT2T3(data.packetNumber, data.time*1000.0)
  })
  //

  $('#videoTest').attr('src', 'test.mp4')
}

function transitionToSelectSla(response) {
  let latencies = response.split(' ');
  console.log(latencies);

  $('.latency-label')[0].innerHTML = `${Math.ceil(parseFloat(latencies[0]) / 10) * 10 + 20} ms`
  $('.latency-label')[1].innerHTML = `${Math.ceil(parseFloat(latencies[1]) / 10) * 10 + 30} ms`
  // $('.latency-label')[2].innerHTML = `${Math.ceil(parseFloat(latencies[2]) / 10) * 10 + 40} ms`

  $('.selection-page').animate({
    left: "-100vw",
    opacity: 0
  }, 300);
  $('.SLA-page').animate({
    left: "0",
    opacity: 1
  }, 300);
}

function sendStopTest() {
  makeCORSRequest('POST', `http://${CLIENT_ADDR}/stoptest`, () => {
    console.log()
  })
  makeCORSRequest('POST', `http://${CLIENT_ADDR}/resettest`)
}
