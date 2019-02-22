// Imports
import React, { Component } from 'react';
import { View } from 'react-native';
import { Provider } from 'react-redux';
import Header from '../components/front-end/header';
import Get_Room from '../components/front-end/getRoom';
import { createStore } from 'redux';
import reducers from '../redux/reducers';


//Create the store
const store = createStore(reducers, {roomName: ''})



export default class Home extends Component {
  render(){
    return (
      <Provider store={store}>
        <View>
          <Header />
          <Get_Room/>

        </View>
      </Provider>

    )
  }
}
