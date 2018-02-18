import { StackNavigator } from 'react-navigation';

import MainTabNavigator from './MainTabNavigator';
import LoginScreen from '../screens/LoginScreen';

export default StackNavigator(
  {
    Main: {
      screen: MainTabNavigator,
    },
    Login: {
      screen: LoginScreen,
    },
  },
  {
    navigationOptions: () => ({
      headerTitleStyle: {
        fontWeight: 'normal',
      },
    }),
  }
);
