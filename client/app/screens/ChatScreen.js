import React from 'react';
import { StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import {
  Spinner,
  View,
  Button
} from 'native-base';
import { GiftedChat, Bubble } from 'react-native-gifted-chat';

import { PTSansText } from '../components/StyledText';
import Colors from '../constants/Colors';
import { sendFirebaseMessage } from '../firebase';
import {
  registerFetchNewMessageAsync,
  fetchEarlierMessagesAsync
} from '../actions/chatActions';

const messagesToLoad = 50;

class ChatScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `Chat with ${navigation.state.params.name}`,
  });

  onSend(messages = []) {
    const _id = this.props.navigation.state.params._id;

    sendFirebaseMessage(_id, messages);
  }

  onLoadEarlier() {
    const {
      messages,
      navigation,
      fetchEarlierMessagesAsync,
    } = this.props;
    const _id = navigation.state.params._id;

    fetchEarlierMessagesAsync(_id, messagesToLoad, messages[messages.length - 1]._id);
  }

  renderBubble (props) {
    return (
      <Bubble
        {...props}
        wrapperStyle={{
          right: {
            backgroundColor: Colors.tintColor
          }
        }}
      />
    )
  }

  render() {
    const { navigation, userProfile, messages } = this.props;
    return (
      <View style={styles.container}>
        <GiftedChat
          messages={messages}
          onSend={messages => this.onSend(messages)}
          messageIdGenerator={() => null}
          user={{
            _id: userProfile.id,
            name: userProfile.formattedName,
            avatar: userProfile.pictureUrl
          }}
          loadEarlier={messages && messages[messages.length - 1]._id !== 0}
          onLoadEarlier={() => this.onLoadEarlier()}
          renderLoading={() => <Spinner color={Colors.tintColor} />}
          renderBubble={this.renderBubble}
          inverted={true}
        />
      </View>
    );
  }

  componentDidMount() {
    const {
      messages,
      navigation,
      registerFetchNewMessageAsync,
      fetchEarlierMessagesAsync
    } = this.props;

    if (!messages) {
      fetchEarlierMessagesAsync(navigation.state.params._id, messagesToLoad);
      registerFetchNewMessageAsync(navigation.state.params._id);
    }
  }
};

const mapStateToProps = (state, ownProps) => ({
  userProfile: state.userProfile,
  messages: state.chat.messageThreads[ownProps.navigation.state.params._id]
});

const mapDispatchToProps = {
  registerFetchNewMessageAsync,
  fetchEarlierMessagesAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ChatScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  }
});
