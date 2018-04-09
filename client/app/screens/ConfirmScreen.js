import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Container,
  View,
  Button,
  Icon
} from 'native-base';
import { connect } from 'react-redux';

import Profile from '../components/Profile';
import { serverURI } from '../constants/Common';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { StyledToast } from '../helpers';

class ConfirmScreen extends React.Component {
  static navigationOptions = {
    title: 'Confirm MeetOver',
  };

  async _sendDecision(decision) {
    const { navigation, accessToken, userId } = this.props;

    const uri = serverURI + `/meetover/decision/${navigation.state.params.profile.id}`;
    const init = {
      method: 'POST',
      body: JSON.stringify({
        _id: navigation.state.params._id,
        status: decision
      }),
      headers: new Headers({
        'Token': accessToken,
        'Identity': userId
      })
    };

    return await fetch(uri, init)
      .catch(err => console.error(err));
  }

  _acceptMeetover() {
    const { navigation } = this.props;
    this._sendDecision('accepted');
    navigation.replace('ChatScreen', navigation.state.params);
  }

  _declineMeetover() {
    const { navigation } = this.props;
    this._sendDecision('declined');
    StyledToast({ text: `Declined ${navigation.state.params.profile.formattedName}'s request` });
    this.props.navigation.goBack();
  }

  render() {
    const { profile } = this.props.navigation.state.params;
    return (
      <Container style={styles.container}>
        {profile && <Profile profile={profile} />}
        <View style={{ flexDirection: 'row' }}>
          <Button
            full
            style={styles.declineButton}
            onPress={() => this._declineMeetover()}
          >
            <PTSansText style={styles.request}>Decline</PTSansText>
          </Button>
          <View style={styles.separator} />
          <Button
            full
            style={styles.acceptButton}
            onPress={() => this._acceptMeetover()}
          >
            <PTSansText style={styles.request}>Accept</PTSansText>
          </Button>
        </View>
      </Container>
    );
  }
}

const mapStateToProps = state => ({
  accessToken: state.userProfile.token.access_token,
  userId: state.userProfile.id
});

export default connect(
  mapStateToProps
)(ConfirmScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  declineButton: {
    flex: 4,
    backgroundColor: Colors.tintColor
  },
  separator: {
    borderRightColor: '#fff',
    borderRightWidth: StyleSheet.hairlineWidth,
  },
  acceptButton: {
    flex: 6,
    backgroundColor: Colors.tintColor
  }
});
