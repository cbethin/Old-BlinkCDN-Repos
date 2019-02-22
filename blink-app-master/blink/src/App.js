// Imports
import React from 'react';
import { StackNavigator } from 'react-navigation';
import Home from './screens/Home'
import Chat from './screens/Chat'


//Component
const App = StackNavigator({
  Home: { screen: Home },
  Chat: { screen: Chat }
})

console.log(StackNavigator)

export default App;
