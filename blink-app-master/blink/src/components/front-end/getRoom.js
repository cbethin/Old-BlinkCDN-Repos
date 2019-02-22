import React, { Component } from 'react'
import { View, Text, TouchableOpacity, TextInput, StyleSheet } from 'react-native'
import { connect } from 'react-redux';
import ROOM_NAME from '../../redux/actions'
import { dispatch } from 'redux'
import { StackNavigator } from 'react-navigation';

class Get_Room extends Component {
   state = {
      roomName: ''
   }
   handleRoomName = (text) => {
      this.setState({ roomName: text })
   }


   render(){
      return (
         <View style = {styles.container}>
            <TextInput style = {styles.input}
               underlineColorAndroid = "transparent"
               placeholder = "Enter Room Name"
               placeholderTextColor = "#076d93"
               autoCapitalize = "none"
               onChangeText = {this.handleRoomName}/>

            <TouchableOpacity
               style = {styles.submitButton}
               onPress = {
                 //Error Here with dispatch function
                 ()=> this.store.dispatch(ROOM_NAME(this.state.roomName))//console.log(this.state.roomName) //this.store.dispatch(ROOM_NAME(this.state.roomName)) //console.log(this.state.roomName)
              }>
            <Text style = {styles.submitButtonText}> Join! </Text>
            </TouchableOpacity>
         </View>
      )
   }
}



const mapStateToProps = (state) => {
return {
  roomName : state.roomName
        }
   }

  

const mapDispatchToProps = () => {
     return {
       ROOM_NAME: ROOM_NAME
     };
};


export default connect(mapStateToProps, mapDispatchToProps ) (Get_Room);

const styles = {
   container: {
      paddingTop: '25%',
      backgroundColor: '#fff'
   },
   // Place Holder Text and Color is above in Text Input Tag
   input: {
      margin: 15,
      height: 40,
      borderRadius: 10,
      borderColor: '#076d93',
      borderWidth: 1
   },
   submitButton: {
      backgroundColor: '#42a5bc',
      padding: 10,
      margin: 15,
      height: 40,
   },
   submitButtonText:{
      color: 'white',
      alignSelf: 'center'
   }
}
