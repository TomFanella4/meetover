import React from 'react';
import { StyleSheet } from 'react-native';
import { Container, Content, View, List, ListItem, Spinner } from 'native-base';
import { connect } from 'react-redux';

import { PTSansText } from '../components/StyledText';
import Colors from '../constants/Colors';
import { chatScreenStrings } from '../constants/Strings';
import IsSearchingBar from '../components/IsSearchingBar';

class ChatsScreen extends React.Component {
  static navigationOptions = {
    title: 'Chats',
  };

  _viewChatThread(_id, name) {
    this.props.navigation.navigate('ChatScreen', { _id, name });
  }

  render() {
    const { threadList } = this.props;

    const threadItems = threadList ?
      threadList.map(thread => (
        <ListItem
          key={thread.name}
          onPress={() => this._viewChatThread(thread._id, thread.name)}
        >
          <PTSansText>{thread.name}</PTSansText>
        </ListItem>
      ))
    :
      null;

    return (
      <Container style={styles.container}>
        <IsSearchingBar />
        {
          threadItems ?
            threadItems.length > 0 ?
              <Content>
                <List>
                  {threadItems}
                </List>
              </Content>
            :
              <View style={styles.noChatsText}>
                <PTSansText>{chatScreenStrings.noChatsText}</PTSansText>
              </View>
          :
            <Spinner color={Colors.tintColor} />
        }
      </Container>
    );
  }
};

const mapStateToProps = state => ({
  threadList: state.chat.threadList
});

export default connect(
  mapStateToProps
)(ChatsScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  noChatsText: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center'
  }
});
