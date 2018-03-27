import React from 'react';
import { Platform, BackHandler } from "react-native";
// import { Notifications } from 'expo';
import { StackNavigator, SwitchNavigator } from 'react-navigation';
import { Root } from 'native-base';

import LoginScreen from '../screens/LoginScreen';
import MainTabNavigator from './MainTabNavigator';
import CreateProfileScreen from '../screens/CreateProfileScreen';
import SettingsScreen from '../screens/SettingsScreen';
import ProfileScreen from '../screens/ProfileScreen';
import Colors from '../constants/Colors';
// import registerForPushNotificationsAsync from '../api/registerForPushNotificationsAsync';

const AppNavigator = StackNavigator(
  {
    Main: MainTabNavigator,
    SettingsScreen: SettingsScreen,
    Profile: ProfileScreen,
  },
  {
    navigationOptions: {
      headerTitleStyle: {
        fontWeight: 'normal',
      },
      headerStyle: {
        backgroundColor: Colors.tintColor,
        marginTop: Platform.OS === 'android' ? -24 : 0,
      },
      headerTintColor: Colors.header,
    },
  }
);

const AuthNavigator = ({ initialRouteName }) => {
  const Nav = SwitchNavigator(
    {
      Login: LoginScreen,
      CreateProfile: CreateProfileScreen,
      App: AppNavigator
    },
    {
      initialRouteName
    }
  );

  return <Nav />
}

export default class AppNavigation extends React.Component {

  componentDidMount() {
    // this._notificationSubscription = this._registerForPushNotifications();
  }

  componentWillUnmount() {
    // this._notificationSubscription && this._notificationSubscription.remove();
  }

  render() {
    const { id, isAuthenticated } = this.props;

    let initialRouteName;
    if (isAuthenticated) {
      initialRouteName = 'App';
    } else if (id) {
      initialRouteName = 'CreateProfile';
    } else {
      initialRouteName = 'Login';
    }

    return (
      <Root>
        <AuthNavigator initialRouteName={initialRouteName} />
      </Root>
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
