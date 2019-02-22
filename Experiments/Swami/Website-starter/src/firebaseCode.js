
var config = {
  apiKey: "AIzaSyA3zZf7U_-on3Xk2wFGsCvZnKoiuLpqKAA",
  authDomain: "swami-database.firebaseapp.com",
  databaseURL: "https://swami-database.firebaseio.com",
  projectId: "swami-database",
  storageBucket: "swami-database.appspot.com",
  messagingSenderId: "333315656264",
};

firebase.initializeApp(config);

var firestore = firebase.firestore();
const settings = {
  timestampsInSnapshots: true
};
firestore.settings(settings);

const docRefTotal = firestore.doc("Latency_Data/t1_t3");
const docRefFirst = firestore.doc("Latency_Data/t1_t2");
const docRefSecond = firestore.doc("Latency_Data/t2_t3");
const outputHeader = document.querySelector("#output");

function getRealTimeUpdatesTotal(callback) {
  docRefTotal.get();
  docRefTotal.onSnapshot((doc) => {
    if (doc && doc.exists) {
      const myData = doc.data();
      // console.log(myData);
      callback(myData);
    }
  });
}

function getRealTimeUpdatesFirst(callback) {
  docRefFirst.get();
  docRefFirst.onSnapshot((doc) => {
    if (doc && doc.exists) {
      const myData = doc.data();
      // console.log(myData);
      callback(myData);
    }
  });
}

function getRealTimeUpdatesSecond(callback) {
  docRefSecond.get();
  docRefSecond.onSnapshot((doc) => {
    if (doc && doc.exists) {
      const myData = doc.data();
      // console.log(myData);
      callback(myData);
    }
  });
}

export {getRealTimeUpdatesTotal, getRealTimeUpdatesFirst, getRealTimeUpdatesSecond}
