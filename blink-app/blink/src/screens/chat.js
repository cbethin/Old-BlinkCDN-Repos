// Imports
import React, { Component } from 'react';
import { View } from 'react-native';
import { Provider } from 'react-redux';
import Header from '../components/front-end/header';


export default class Chat extends Component {
  render(){
    return (
      <Provider store={store}>
        <View>
          <Header />

        </View>
      </Provider>

    )
  }
}
