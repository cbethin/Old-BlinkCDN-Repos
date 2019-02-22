// Create the XHR object.
function createCORSRequest(method, url) {
  var xhr = new XMLHttpRequest();
  if ("withCredentials" in xhr) {
    // XHR for Chrome/Firefox/Opera/Safari.
    xhr.open(method, url, true);
  } else if (typeof XDomainRequest != "undefined") {
    // XDomainRequest for IE.
    xhr = new XDomainRequest();
    xhr.open(method, url);
  } else {
    // CORS not supported.
    xhr = null;
  }

  return xhr;
}

// Send CORS get request
function makeCORSRequest(method, url, callback) {
  var xhr = createCORSRequest(method, url);

  if (!xhr) {
    console.log('CORS not supported');
    return;
  }

  xhr.onload = function() {
    // console.log(xhr.responseText);
    callback(xhr.response)
  }

  xhr.onerror = () => {
    console.log("ERROR: Whoops, error making request");
  }

  xhr.send();
}

export {makeCORSRequest}
