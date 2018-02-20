import React from 'react';
import { StyleSheet } from 'react-native';
import { NavigationActions } from 'react-navigation';
import { AuthSession } from 'expo';
import { Button, View, Spinner } from 'native-base';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { LI_APP_ID } from 'react-native-dotenv';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { saveProfileAndLoginAsync } from '../actions';

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

      const uri = `https://meetover.herokuapp.com/login/${result.params.code}`;
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      const userProfile = await response.json();
      userProfile = { ...userProfile, isAuthenticated: true }

      this.props.saveProfileAndLoginAsync(userProfile)
      .then(() => this.setState({ isLoading: false }));
    }
  }
}

const mapDispatchToProps = {
  saveProfileAndLoginAsync
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
