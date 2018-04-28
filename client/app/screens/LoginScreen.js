import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Button,
  Form,
  Input,
  Item,
  Spinner,
  View,
} from 'native-base';

import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { authenticateAndCreateProfile, imitateUser } from '../actions/userActions';
import { allowImitate } from '../constants/Common';

class LoginScreen extends React.Component {
  static navigationOptions = {
    header: null
  };

  state = {
    isLoading: false,
    imitateUID: '',
  }

  _renderImitateForm() {
    if (allowImitate) {
      return (
        <Form>
          <Item>
            <Input
              placeholder='uid'
              onChangeText={imitateUID => this.setState({ imitateUID })}
            />
            <Button
              style={styles.button}
              onPress={this._handleImitatePressAsync}
            >
              <PTSansText>Imitate</PTSansText>
            </Button>
          </Item>
        </Form>
      );
    }

    return null;
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
        {this._renderImitateForm()}
      </View>
    );
  }

  _handleImitatePressAsync = async () => {
    const { imitateUID } = this.state;
    const { imitateUser, navigation } = this.props;

    this.setState({ isLoading: true });
    await imitateUser(imitateUID);
    this.setState({ isLoading: false });

    navigation.navigate('App');
  }

  _handleSignInPressAsync = async () => {
    const { authenticateAndCreateProfile, navigation } = this.props;
    this.setState({ isLoading: true });
    const result = await authenticateAndCreateProfile();
    this.setState({ isLoading: false });

    if (result.type === 'success') {
      if (result.isAuthenticated) {
        navigation.navigate('App');
      } else {
        navigation.navigate('CreateProfile');
      }
    }
  }
}

const mapDispatchToProps = {
  authenticateAndCreateProfile,
  imitateUser
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
