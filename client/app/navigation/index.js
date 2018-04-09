import React from 'react';
import { Platform, BackHandler } from "react-native";
import { StackNavigator, SwitchNavigator } from 'react-navigation';
import { Root } from 'native-base';

import LoginScreen from '../screens/LoginScreen';
import MainTabNavigator from './MainTabNavigator';
import CreateProfileScreen from '../screens/CreateProfileScreen';
import SettingsScreen from '../screens/SettingsScreen';
import ProfileScreen from '../screens/ProfileScreen';
import RequestScreen from '../screens/RequestScreen';
import ConfirmScreen from '../screens/ConfirmScreen';
import ChatScreen from '../screens/ChatScreen';
import Colors from '../constants/Colors';

const AppNavigator = StackNavigator(
  {
    Main: MainTabNavigator,
    SettingsScreen: SettingsScreen,
    ProfileScreen: ProfileScreen,
    RequestScreen: RequestScreen,
    ConfirmScreen: ConfirmScreen,
    ChatScreen: ChatScreen,
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

export default Navigation = ({ id, isAuthenticated }) => {

  let initialRouteName;
  if (isAuthenticated) {
    initialRouteName = 'App';
  } else if (id) {
    initialRouteName = 'CreateProfile';
  } else {
    initialRouteName = 'Login';
  }

  const AuthNavigator = SwitchNavigator(
    {
      Login: LoginScreen,
      CreateProfile: CreateProfileScreen,
      App: AppNavigator
    },
    {
      initialRouteName
    }
  );

  return (
    <Root>
      <AuthNavigator />
    </Root>
  );
};
