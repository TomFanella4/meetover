import React from 'react';
import { Platform, StatusBar, StyleSheet, View, YellowBox } from 'react-native';
import { Root } from "native-base";
import Expo, { AppLoading, Asset, Font } from 'expo';
import { Ionicons } from '@expo/vector-icons';
import AppNavigation from './navigation';

import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';

import reducers from './reducers';
import { signInToFirebase } from './firebase';

YellowBox.ignoreWarnings(['Setting a timer']);

export default class App extends React.Component {
  state = {
    isLoadingComplete: false,
    store: {},
    userProfile: {}
  };

  render() {
    if (!this.state.isLoadingComplete && !this.props.skipLoadingScreen) {
      return (
        <AppLoading
          startAsync={this._loadAll}
          onError={this._handleLoadingError}
          onFinish={this._handleFinishLoading}
        />
      );
    } else {
      return (
        <Provider store={this.state.store}>
          <View style={styles.container}>
            {Platform.OS === 'ios' && <StatusBar barStyle="light-content" />}
            {/* {Platform.OS === 'android' && <View style={styles.statusBarUnderlay} />} */}
            <AppNavigation
              id={this.state.userProfile.id}
              isAuthenticated={this.state.userProfile.isAuthenticated}
            />
          </View>
        </Provider>
      );
    }
  }

  _loadAll = async () => {
    const userProfile = (await this._loadResourcesAsync())[2];
    userProfile && await this._validateCustomToken(userProfile);
  };

  _loadResourcesAsync = async () => {
    return Promise.all([
      Asset.loadAsync([
        require('../assets/images/current_location.png'),
        require('../assets/images/icon.png'),
      ]),
      Font.loadAsync({
        // This is the font that we are using for our tab bar
        ...Ionicons.font,
        // PT-Sans is the main font for MeetOver
        'pt-sans': require('../assets/fonts/PT_Sans-Web-Regular.ttf'),
        'Roboto': require("native-base/Fonts/Roboto.ttf"),
        'Roboto_medium': require("native-base/Fonts/Roboto_medium.ttf")
      }),
      Expo.SecureStore.getItemAsync('userProfile')
      .then(userProfileString => {
        if (userProfileString) {
          return JSON.parse(userProfileString);
        }
      })
    ]);
  };

  _validateCustomToken = async userProfile => {
    const { firebaseCustomToken, token, id } = userProfile;

    const verifiedToken = await signInToFirebase(firebaseCustomToken, token.access_token, id)
    .catch(() => {
      Expo.SecureStore.deleteItemAsync('userProfile');
      throw 'Unable to verify Firebase custom token.';
    });

    userProfile.firebaseCustomToken = verifiedToken;

    this.setState({ userProfile });
    Expo.SecureStore.setItemAsync('userProfile', JSON.stringify(userProfile));
  };

  _handleLoadingError = error => {
    // In this case, you might want to report the error to your error
    // reporting service, for example Sentry
    console.warn(error);
  };

  _handleFinishLoading = () => {
    this.setState({
      isLoadingComplete: true,
      store: createStore(
        reducers,
        { userProfile: this.state.userProfile },
        applyMiddleware(thunk)
      )
    });
  };
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  statusBarUnderlay: {
    height: 24,
    backgroundColor: 'rgba(0,0,0,0.2)',
  },
});
