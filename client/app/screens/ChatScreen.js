import React from 'react';
import { StyleSheet } from 'react-native';
import { GiftedChat } from 'react-native-gifted-chat'
import * as firebase from 'firebase';

import { PTSansText } from '../components/StyledText'

export default class ChatScreen extends React.Component {
  static navigationOptions = {
    title: 'Chat',
  };

  state = {
    messages: [],
  }

  componentWillMount() {
    // this.setState({
    //   messages: [
    //     {
    //       _id: 1,
    //       text: 'Hello developer',
    //       createdAt: new Date(),
    //       user: {
    //         _id: 2,
    //         name: 'React Native',
    //         avatar: 'https://i.imgur.com/wq43v5T.jpg',
    //       },
    //     },
    //   ],
    // })
    // console.log(this.state.messages);

    var userId = firebase.auth().currentUser.uid;
    var messages = firebase.database().ref('messages/' + userId);
    messages.on('value', snapshot => {
      console.log(snapshot.val());
      snapshot.val() &&
      this.setState(previousState => ({
        messages: snapshot.val().messages
      }))
    })
    // .then(snapshot => {
    //   console.log(snapshot.val());
    // })
    // .catch(err => console.error(err));
  }

  onSend(messages = []) {
    // this.setState(previousState => ({
    //   messages: GiftedChat.append(previousState.messages, messages),
    // }))
    var userId = firebase.auth().currentUser.uid;
    firebase.database().ref('messages/' + userId).set({
      messages: GiftedChat.append(this.state.messages, messages),
      // messages: [
      //     {
      //       _id: 1,
      //       text: 'Hello developer',
      //       createdAt: JSON.stringify(new Date()),
      //       user: {
      //         _id: 2,
      //         name: 'React Native',
      //         avatar: 'https://i.imgur.com/wq43v5T.jpg',
      //       },
      //     },
      //   ],
    });
  }

  render() {
    // console.log(this.state.messages);
    return (
      <GiftedChat
        messages={this.state.messages}
        onSend={messages => this.onSend(messages)}
        user={{
          _id: 1,
        }}
      />
    );
  }
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 15,
    backgroundColor: '#fff',
  }
});
