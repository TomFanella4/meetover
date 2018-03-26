import React from 'react';
import { StyleSheet } from 'react-native';
import { Button, View, Spinner } from 'native-base';

import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { authenticateAndCreateProfile } from '../actions/userActions';

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
    const { authenticateAndCreateProfile, navigation } = this.props;
    this.setState({ isLoading: true });
    const resultType = await authenticateAndCreateProfile();
    this.setState({ isLoading: false });
    resultType === 'success' && navigation.navigate('CreateProfile');
  }
}

const mapDispatchToProps = {
  authenticateAndCreateProfile
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
