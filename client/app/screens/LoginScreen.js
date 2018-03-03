import React from 'react';
import { StyleSheet } from 'react-native';
import { AuthSession } from 'expo';
import { Button, View, Spinner } from 'native-base';

import { connect } from 'react-redux';

import { LI_APP_ID } from 'react-native-dotenv';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { createProfile, saveProfileAndLoginAsync } from '../actions';
import { fetchIdToken } from '../firebase';

class LoginScreen extends React.Component {
  static navigationOptions = {
    header: null
  };

  state = {
    isLoading: false
  }

  render() {
    return (
      <View style={styles.container}>
        <View style={styles.title}>
          <PTSansText style={styles.titleText}>
            MeetOver
          </PTSansText>
          <PTSansText style={styles.subtitleText}>
            Connecting Professionals on the Fly
          </PTSansText>
        </View>
        <View style={styles.content}>
          {
            this.state.isLoading ?
              <Spinner color={Colors.tintColor} />
            :
              <Button
                style={styles.button}
                onPress={this._handleSignInPressAsync}
              >
                <PTSansText>Sign in with LinkedIn</PTSansText>
              </Button>
          }
          <PTSansText style={styles.termsText}>
            By signing in you accept our terms of service
          </PTSansText>
        </View>
      </View>
    );
  }

  _handleSignInPressAsync = async () => {
    const redirectUri = AuthSession.getRedirectUrl();
    const result = await AuthSession.startAsync({
      authUrl:
        `https://www.linkedin.com/oauth/v2/authorization?response_type=code` +
        `&client_id=${LI_APP_ID}` +
        `&redirect_uri=${encodeURIComponent(redirectUri)}` +
        `&state=meetover_testing`
    });

    if (result.type === 'success') {
      this.setState({ isLoading: true });

      const uri = `https://meetover.herokuapp.com/login/${result.params.code}` +
        `?redirect_uri=${encodeURIComponent(redirectUri)}`;
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      const { profile, token, firebaseCustomToken } = await response.json();
      const firebaseIdToken = await fetchIdToken(firebaseCustomToken)
        .catch(err => null);

      this.props.createProfile({ ...profile, token, firebaseCustomToken });
      this.setState({ isLoading: false });
    }
  }
}

const mapDispatchToProps = {
  createProfile
};

export default connect(
  null,
  mapDispatchToProps
)(LoginScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  title: {
    backgroundColor: Colors.tintColor,
    flex: 5,
    alignItems: 'center',
    justifyContent: 'center',
  },
  content: {
    flex: 5,
    alignItems: 'center',
    justifyContent: 'center',
  },
  titleText: {
    color: 'white',
    fontSize: 70
  },
  subtitleText: {
    color: 'white',
    fontSize: 20
  },
  termsText: {
    color: 'grey',
    fontSize: 10
  },
  button: {
    backgroundColor: Colors.tintColor,
    alignSelf: 'center'
  }
});
