import React from 'react';
import { connect } from 'react-redux';
import { Platform, BackHandler } from "react-native";
// import { Notifications } from 'expo';
import { addNavigationHelpers, StackNavigator, NavigationActions } from 'react-navigation';

import LoginScreen from '../screens/LoginScreen';
import MainTabNavigator from './MainTabNavigator';
import CreateProfileScreen from '../screens/CreateProfileScreen';
import SettingsScreen from '../screens/SettingsScreen';
import ProfileScreen from '../screens/ProfileScreen';
import { addListener } from '../store/middleware';
import Colors from '../constants/Colors';
// import registerForPushNotificationsAsync from '../api/registerForPushNotificationsAsync';

const defaultNavigationOptions = {
  headerTitleStyle: {
    fontWeight: 'normal',
  },
  headerStyle: {
    backgroundColor: Colors.tintColor,
    marginTop: Platform.OS === 'android' ? -24 : 0,
  },
  headerTintColor: Colors.header,
};

export const AppNavigator = StackNavigator(
  {
    Main: {
      screen: MainTabNavigator,
    },
    Login: {
      screen: LoginScreen,
    },
    CreateProfile: {
      screen: CreateProfileScreen,
    },
    SettingsScreen: {
      screen: SettingsScreen,
    },
    Profile: {
      screen: ProfileScreen,
      navigationOptions: ({ navigation }) => ({
        ...defaultNavigationOptions,
        title: `${navigation.state.params.name}'s Profile`,
      })
    },
  },
  {
    navigationOptions: defaultNavigationOptions,
  }
);

class AppNavigation extends React.Component {
  componentDidMount() {
    BackHandler.addEventListener("hardwareBackPress", this.onBackPress);
    // this._notificationSubscription = this._registerForPushNotifications();
  }

  componentWillUnmount() {
    BackHandler.removeEventListener("hardwareBackPress", this.onBackPress);
    // this._notificationSubscription && this._notificationSubscription.remove();
  }

  onBackPress = () => {
    const { dispatch, navigationState } = this.props;
    if (navigationState.stateForLoggedIn.index <= 1) {
      BackHandler.exitApp();
      return;
    }
    dispatch(NavigationActions.back());
    return true;
  };


  render() {
    const { navigationState, dispatch, isAuthenticated } = this.props;
    const state = isAuthenticated
      ? navigationState.stateForLoggedIn
      : navigationState.stateForLoggedOut;

    return (
      <AppNavigator navigation={addNavigationHelpers({ dispatch, state, addListener })} />
    );
  }

  // _registerForPushNotifications() {
  //   // Send our push token over to our backend so we can receive notifications
  //   // You can comment the following line out if you want to stop receiving
  //   // a notification every time you open the app. Check out the source
  //   // for this function in api/registerForPushNotificationsAsync.js
  //   registerForPushNotificationsAsync();
  //
  //   // Watch for incoming notifications
  //   this._notificationSubscription = Notifications.addListener(this._handleNotification);
  // }

  // _handleNotification = ({ origin, data }) => {
  //   console.log(`Push notification ${origin} with data: ${JSON.stringify(data)}`);
  // };
}

const mapStateToProps = state => ({
  navigationState: state.navigation,
  isAuthenticated: state.userProfile.isAuthenticated
});

export default connect(
  mapStateToProps
)(AppNavigation);
