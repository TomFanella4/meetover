import React from 'react';
import { ScrollView, StyleSheet } from 'react-native';

import { PTSansText } from '../components/StyledText'

export default class ChatsScreen extends React.Component {
  static navigationOptions = {
    title: 'Chats',
  };

  render() {
    return (
      <ScrollView style={styles.container}>
        <PTSansText>Chats Screen</PTSansText>
      </ScrollView>
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
