import React from 'react';
import { StyleSheet } from 'react-native';
import {
  View,
  Button,
  Icon,
  Container,
  Spinner
} from 'native-base';
import { connect } from 'react-redux';
import { find } from 'lodash';

import Profile from '../components/Profile';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { separator, serverURI } from '../constants/Common';
import { StyledToast } from '../helpers';

class RequestScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `${navigation.state.params.profile.formattedName}'s Profile`,
  });

  state = {
    buttonDisabled: true,
    meetOverLoading: false
  };

  componentDidMount() {
    const { threadList, navigation } = this.props;

    if (threadList !== null) {
      // Thread list has been fetched from Firebase
      this.setState({ buttonDisabled: false });
    }
    this._blurSub = navigation.addListener('didBlur', () => this._onBlur());
  }

  componentDidUpdate(prevProps) {
    const { threadList } = this.props;
    const prevThreadList = prevProps.threadList;

    if (prevThreadList === null && threadList !== null) {
      // Thread list has been fetched from Firebase
      this.setState({ buttonDisabled: false });
    }
  }

  _onBlur() {
    this.setState({ buttonDisabled: false, meetOverLoading: false });
  }

  _renderLoading() {
    return (
      <Container style={styles.container}>
        <View style={styles.loadingView}>
          <PTSansText>Sending MeetOver Request...</PTSansText>
          <Spinner color={Colors.tintColor} />
        </View>
      </Container>
    );
  }

  async _initiateMeetover() {
    const { navigation, signedInProfile, threadList } = this.props;
    const { formattedName, id } = navigation.state.params.profile;
    const signedInId = signedInProfile.id;
    const accessToken = signedInProfile.token.access_token;
    let threadId;

    this.setState({ buttonDisabled: true, meetOverLoading: true });

    if (signedInId < id) {
      threadId = signedInId + separator + id;
    } else {
      threadId = id + separator + signedInId;
    }

    const exists = (find(threadList, { '_id': threadId }) !== undefined);

    if (!exists) {
      const uri = `${serverURI}/meetover/${id}`;
      const init = {
        method: 'POST',
        headers: new Headers({
          'Token': accessToken,
          'Identity': signedInId
        })
      };

      const response = await fetch(uri, init)
        .catch(err => console.log(err));

      if (response.status !== 200) {
        console.log('Could not initiate meetover');
        console.log(response);

        StyledToast({
          text: 'Could not initate MeetOver',
          buttonText: 'Okay',
          type: 'danger',
          duration: 3000,
        });

        return;
      }
    }

    navigation.navigate('ChatScreen', { _id: threadId, name: formattedName });
  };

  render() {
    const { profile } = this.props.navigation.state.params;
    const { buttonDisabled, meetOverLoading } = this.state;

    if(meetOverLoading) {
      return this._renderLoading();
    } else {
      return (
        <Container style={styles.container}>
          <Profile profile={profile} />
          <Button
            iconLeft
            full
            style={!buttonDisabled ? styles.chatButton : null}
            disabled={buttonDisabled}
            onPress={() => this._initiateMeetover()}
          >
            <Icon name='chatboxes' />
            <PTSansText style={styles.request}>Request MeetOver</PTSansText>
          </Button>
        </Container>
      );
    }
  }

  componentWillUnmount() {
    this._blurSub.remove();
  }
};

const mapStateToProps = state => ({
  signedInProfile: state.userProfile,
  threadList: state.chat.threadList
});

export default connect(
  mapStateToProps
)(RequestScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  loadingView: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center'
  },
  chatButton: {
    backgroundColor: Colors.tintColor
  },
  request: {
    fontSize: 18
  }
});
