  // Initialize Firebase
  var config = {
    apiKey: "AIzaSyA3zZf7U_-on3Xk2wFGsCvZnKoiuLpqKAA",
    authDomain: "swami-database.firebaseapp.com",
    databaseURL: "https://swami-database.firebaseio.com",
    projectId: "swami-database",
    storageBucket: "swami-database.appspot.com",
    messagingSenderId: "333315656264"
  };
  firebase.initializeApp(config);

var firestore = firebase.firestore();

const docRef = firestore.doc("sampleData/inspiration")
const outputHeader = document.querySelector("#output")

function getRealTimeUpdates() {
  docRef.onSnapshot(function (doc) {
    if (doc && doc.exists) {
      const myData = doc.data();
      outputHeader.innerText = "Status" + myData.Latency
    }
  }
}

getRealTimeUpdates()
