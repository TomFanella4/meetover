import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Container,
  Content,
  View,
  List,
  ListItem,
  Left,
  Body,
  Spinner
} from 'native-base';
import { connect } from 'react-redux';

import { PTSansText } from '../components/StyledText';
import { ProfileImage } from '../components/ProfileImage';
import Colors from '../constants/Colors';
import { chatScreenStrings } from '../constants/Strings';
import IsSearchingBar from '../components/IsSearchingBar';

class ChatsScreen extends React.Component {
  static navigationOptions = {
    title: 'Chats',
  };

  _viewChatThread(thread) {
    thread.status === 'pending' && thread.origin === 'receiver' ?
      this.props.navigation.navigate('ConfirmScreen', thread)
    :
      this.props.navigation.navigate('ChatScreen', thread);
  }

  render() {
    const { threadList } = this.props;
    const threadItems = [];

    threadList.forEach(thread => (
      thread.status !== 'declined' && threadItems.push(
        <ListItem
          key={thread._id}
          onPress={() => this._viewChatThread(thread)}
          avatar
        >
          <Left>
            <ProfileImage pictureUrl={thread.profile.pictureUrl} />
          </Left>
          <Body>
            <PTSansText style={styles.name}>{thread.profile.formattedName}</PTSansText>
            <PTSansText note>{thread.profile.headline}</PTSansText>
          </Body>
        </ListItem>
      )
    ));

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
  },
  name: {
    fontSize: 20
  }
});
