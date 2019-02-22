  'use strict'
  // Global Parameters from .mpd file
  var file;  // MP4 file
  var type;  // Type of file
  var codecs; //  Codecs allowed
  var width;  //  Native width and height
  var height;

  // Elements
  var videoElement = document.getElementById('myVideo');
  var playButton = document.getElementById("load");
  videoElement.poster = "poster.png";

  // Description of initialization segment, and approx segment lengths
  var initialization;
  var segDuration;
  var vidDuration;

  // Video parameters
  var bandwidth; // bitrate of video

  // Parameters to drive segment loop
  var index = 0; // Segment to get
  var segments;
  var curIndex = document.getElementById("curIndex"); // Playing segment
  var segLength = document.getElementById("segLength");

  // Source and buffers
  var mediaSource;
  var videoSource;

  // Parameters to drive fetch loop
  var segCheck;
  var lastTime = 0;
  var bufferUpdated = false;

  // Flags to keep things going
  var lastMpd = "";
  var vTime = document.getElementById("curTime");
  var requestId = 0;

  // Click event handler for load button
  playButton.addEventListener("click", function () {
    //  If video is paused then check for file change
    if (videoElement.paused == true) {
      // Retrieve mpd file, and set up video
      var curMpd = document.getElementById("filename").value;
      //  If current mpd file is different then last mpd file, load it.
      if (curMpd != lastMpd) {
        //  Cancel display of current video position
        window.cancelAnimationFrame(requestId);
        lastMpd = curMpd;
        getData(curMpd);
      } else {
        //  No change, just play
        videoElement.play();
      }
    } else {
      //  Video was playing, now pause it
      videoElement.pause();
    }
  }, false);

  // Do a little trickery, start video when you click the video element
  videoElement.addEventListener("click", function () {
    playButton.click();
  }, false);

  // Event handler for the video element errors
  document.getElementById("myVideo").addEventListener("error", function (e) {
    log("video error: " + e.message);
  }, false);


  // Gets the mpd file and parses it
  function getData(url) {
    if (url !== "") {
      var xhr = new XMLHttpRequest(); // Set up xhr request
      xhr.open("GET", url, true); // Open the request
      xhr.responseType = "text"; // Set the type of response expected
      xhr.send();

      //  Asynchronously wait for the data to return
      xhr.onreadystatechange = function () {
        if (xhr.readyState == xhr.DONE) {
          var tempoutput = xhr.response;
          var parser = new DOMParser(); //  Create a parser object

          // Create an xml document from the .mpd file for searching
          var xmlData = parser.parseFromString(tempoutput, "text/xml", 0);
          log("parsing mpd file");

          // Get and display the parameters of the .mpd file
          getFileType(xmlData);

          // Set up video object, buffers, etc
          setupVideo();

          // Initialize a few variables on reload
          clearVars();
        }
      }

      // Report errors if they happen during xhr
      xhr.addEventListener("error", function (e) {
        log("Error: " + e + " Could not load url.");
      }, false);
    }
  }

  // Retrieve parameters from our stored .mpd file
  function getFileType(data) {
    try {
      file = data.querySelectorAll("BaseURL")[0].textContent.toString();
      var rep = data.querySelectorAll("Representation");
      type = rep[0].getAttribute("mimeType");
      codecs = rep[0].getAttribute("codecs");
      width = rep[0].getAttribute("width");
      height = rep[0].getAttribute("height");
      bandwidth = rep[0].getAttribute("bandwidth");

      var ini = data.querySelectorAll("Initialization");
      initialization = ini[0].getAttribute("range");
      segments = data.querySelectorAll("SegmentURL");

      // Get the length of the video per the .mpd file
      //   since the video.duration will always say infinity
      var period = data.querySelectorAll("Period");
      var vidTempDuration = period[0].getAttribute("duration");
      vidDuration = parseDuration(vidTempDuration); // display length

      var segList = data.querySelectorAll("SegmentList");
      segDuration = segList[0].getAttribute("duration");

    } catch (er) {
      log(er);
      return;
    }
    showTypes();  // Display parameters
  }

  // Display parameters from the .mpd file
  function showTypes() {
    var display = document.getElementById("myspan");
    var spanData;
    spanData = "<h3>Reported values:</h3><ul><li>Media file: " + file + "</li>";
    spanData += "<li>Type: " + type + "</li>";
    spanData += "<li>Codecs: " + codecs + "</li>";
    spanData += "<li>Width: " + width + " -- Height: " + height + "</li>";
    spanData += "<li>Bandwidth: " + bandwidth + "</li>";
    spanData += "<li>Initialization Range: " + initialization + "</li>";
    spanData += "<li>Segment length: " + segDuration / 1000 + " seconds</li>";
    spanData += "<li>" + vidDuration + "</li>";
    spanData += "</ul>";
    display.innerHTML = spanData;
    document.getElementById("numIndexes").innerHTML = segments.length;
    document.getElementById("curInfo").style.display = "block";
    document.getElementById("curInfo").style.display = "block";
  }


  function render() {
    // Display current video position
    vTime.innerText = formatTime(videoElement.currentTime);
    // Recall this function when available
    requestId = window.requestAnimationFrame(render);
  }

  // Create mediaSource and initialize video
  function setupVideo() {
    clearLog(); // Clear console log

    //  Create the media source
    if (window.MediaSource) {
      mediaSource = new window.MediaSource();
     } else {
      log("mediasource or syntax not supported");
      return;
    }
    var url = URL.createObjectURL(mediaSource);
    videoElement.pause();
    videoElement.src = url;
    videoElement.width = width;
    videoElement.height = height;

    // Wait for event that tells us that our media source object is
    //   ready for a buffer to be added.
    mediaSource.addEventListener('sourceopen', function (e) {
      try {
        videoSource = mediaSource.addSourceBuffer('video/mp4');
        initVideo(initialization, file);
      } catch (e) {
        log('Exception calling addSourceBuffer for video', e);
        return;
      }
    },false);

    // Handler to switch button text to Play
    videoElement.addEventListener("pause", function () {
      playButton.innerText = "Play";
    }, false);

    // Handler to switch button text to pause
    videoElement.addEventListener("playing", function () {
      playButton.innerText = "Pause";
    }, false);
    // Remove the handler for the timeupdate event
    videoElement.addEventListener("ended", function () {
      videoElement.removeEventListener("timeupdate", checkTime);
    }, false);
  }

  //  Load video's initialization segment
  function initVideo(range, url) {
    var xhr = new XMLHttpRequest();
    if (range || url) { // make sure we've got incoming params

      // Set the desired range of bytes we want from the mp4 video file
      xhr.open('GET', url);
      xhr.setRequestHeader("Range", "bytes=" + range);
      segCheck = (timeToDownload(range) * .8).toFixed(3); // use .8 as fudge factor
      xhr.send();
      xhr.responseType = 'arraybuffer';
      try {
        xhr.addEventListener("readystatechange", function () {
           if (xhr.readyState == xhr.DONE) { // wait for video to load
            // Add response to buffer
            try {
              videoSource.appendBuffer(new Uint8Array(xhr.response));
              // Wait for the update complete event before continuing
              videoSource.addEventListener("update",updateFunct, false);

            } catch (e) {
              log('Exception while appending initialization content', e);
            }
          }
        }, false);
      } catch (e) {
        log(e);
      }
    } else {
      return // No value for range or url
    }
  }

  function updateFunct() {
    //  This is a one shot function, when init segment finishes loading,
    //    update the buffer flag, call getStarted, and then remove this event.
    bufferUpdated = true;
    getStarted(file); // Get video playback started
    //  Now that video has started, remove the event listener
    videoSource.removeEventListener("update", updateFunct);
  }

  //  Play our file segments
  function getStarted(url) {

    //  Start by loading the first segment of media
    playSegment(segments[index].getAttribute("mediaRange").toString(), url);

    // Start showing video time
    requestId = window.requestAnimationFrame(render);

    // Display current index
    curIndex.textContent = index + 1;
    index++;

    //  Continue in a loop where approximately every x seconds reload the buffer
    videoElement.addEventListener("timeupdate", fileChecks, false);

  }
  //  Get video segments
  function fileChecks() {
    // If we're ok on the buffer, then continue
    if (bufferUpdated == true) {
      if (index < segments.length) {
        // Loads next segment when time is close to the end of the last loaded segment
        if ((videoElement.currentTime - lastTime) >= segCheck) {
          playSegment(segments[index].getAttribute("mediaRange").toString(), file);
          lastTime = videoElement.currentTime;
          curIndex.textContent = index + 1; // Display current index
          index++;
        }
      } else {
        videoElement.removeEventListener("timeupdate", fileChecks, false);
      }
    }
  }

  //  Play segment plays a byte range (format nnnn-nnnnn) of a media file
  function playSegment(range, url) {
    var xhr = new XMLHttpRequest();
    if (range || url) { // Make sure we've got incoming params
      xhr.open('GET', url);
      xhr.setRequestHeader("Range", "bytes=" + range);
      xhr.send();
      xhr.responseType = 'arraybuffer';
      try {
        xhr.addEventListener("readystatechange", function () {
          if (xhr.readyState == xhr.DONE) { //wait for video to load
            //  Calculate when to get next segment based on time of current one
              segCheck = (timeToDownload(range) * .8).toFixed(3); // Use .8 as fudge factor
              segLength.textContent = segCheck;
            // Add received content to the buffer
            try {
              videoSource.appendBuffer(new Uint8Array(xhr.response));
            } catch (e) {
              log('Exception while appending', e);
            }
          }
        }, false);
      } catch (e) {
        log(e);
        return // No value for range
      }
    }
  }

  //  Logs messages to the console
  function log(s) {
    //  send to console
    //    you can also substitute UI here
    console.log(s);
  };

  //  Clears the log
  function clearLog() {
    console.clear();
  }

  function clearVars() {
    index = 0;
    lastTime = 0;
  }

  function timeToDownload(range) {
    var vidDur = range.split("-");
    // Time = size * 8 / bitrate
    return (((vidDur[1] - vidDur[0]) * 8) / bandwidth)
  }

  // Converts mpd time to human time
  function parseDuration(pt) {
    // Parse time from format "PT#H#M##.##S"
    var ptTemp = pt.split("T")[1];
    ptTemp = ptTemp.split("H")
    var hours = ptTemp[0];
    var minutes = ptTemp[1].split("M")[0];
    var seconds = ptTemp[1].split("M")[1].split("S")[0];
    var hundredths = seconds.split(".");
    //  Display the length of video (taken from .mpd file, since video duration is infinate)
    return "Video length: " + hours + ":" + pZ(minutes, 2) + ":" + pZ(hundredths[0], 2) + "." + hundredths[1];

  }


  //  Converts time in seconds into a string HH:MM:SS.ss
  function formatTime(timeSec) {
    var seconds = timeSec % 60;                                 // Get seconds portion
    var minutes = ((timeSec - seconds) / 60) % 60;              // Get minutes portion
    var hours = ((timeSec - seconds - (minutes * 60))) / 3600;  // Get hours portion
    seconds = seconds.toFixed(2);   // Restrict to 2 places (hundredths of seconds)
    var dispSeconds = seconds.toString().split(".");
    return (pZ(hours, 2) + ":" + pZ(minutes, 2) + ":" + pZ(dispSeconds[0], 2) + "." + pZ(dispSeconds[1], 2));
  }

  //  Pad digits with zeros if needed
  function pZ(value, padCount) {
    var tNum = value + '';
    while (tNum.length < padCount) {
      tNum = "0" + tNum;
    }
    return tNum;
  }
