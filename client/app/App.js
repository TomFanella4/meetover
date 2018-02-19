import React from 'react';
import { Platform, StatusBar, StyleSheet, View } from 'react-native';
import { Root } from "native-base";
import Expo, { AppLoading, Asset, Font } from 'expo';
import { Ionicons } from '@expo/vector-icons';
import AppNavigation from './navigation';

import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';

import { middleware } from './store/middleware';
import reducers from './reducers';

export default class App extends React.Component {
  state = {
    isLoadingComplete: false,
    store: {},
    preloadedState: {}
  };

  render() {
    if (!this.state.isLoadingComplete && !this.props.skipLoadingScreen) {
      return (
        <AppLoading
          startAsync={this._loadResourcesAsync}
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
            <AppNavigation />
          </View>
        </Provider>
      );
    }
  }

  _loadResourcesAsync = async () => {
    return Promise.all([
      // Asset.loadAsync([
      //   require('./assets/images/robot-dev.png'),
      //   require('./assets/images/robot-prod.png'),
      // ]),
      Font.loadAsync({
        // This is the font that we are using for our tab bar
        ...Ionicons.font,
        // PT-Sans is the main font for MeetOver
        'pt-sans': require('../assets/fonts/PT_Sans-Web-Regular.ttf'),
        'Roboto': require("native-base/Fonts/Roboto.ttf"),
        'Roboto_medium': require("native-base/Fonts/Roboto_medium.ttf")
      }),
      Expo.SecureStore.getItemAsync('userProfile')
      .then(userProfile => userProfile && this.setState({
        preloadedState: {
          userProfile: JSON.parse(userProfile)
        }
      }))
    ]);
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
        this.state.preloadedState,
        applyMiddleware(...middleware)
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
