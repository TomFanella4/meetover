import React from 'react';
import { Platform } from 'react-native';
import { Constants, Permissions, Notifications } from 'expo';
import { Ionicons } from '@expo/vector-icons';
import { TabNavigator, TabBarBottom } from 'react-navigation';
import { Button } from 'native-base';
import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { registerFetchThreadListAsync } from '../actions/chatActions';
import { modifyFirebaseUserState } from '../firebase';

import ListScreen from '../screens/ListScreen';
import MapScreen from '../screens/MapScreen';
import ChatsScreen from '../screens/ChatsScreen';

const MainTabNavigator = TabNavigator(
  {
    List: ListScreen,
    Map: MapScreen,
    Chats: ChatsScreen,
  },
  {
    navigationOptions: ({ navigation }) => ({
      tabBarIcon: ({ focused, tintColor }) => {
        const { routeName } = navigation.state;
        let iconName;
        switch (routeName) {
          case 'List':
            iconName =
              Platform.OS === 'ios' ? `ios-list${focused ? '' : '-outline'}` : 'md-list';
            break;
          case 'Map':
            iconName = Platform.OS === 'ios' ? `ios-map${focused ? '' : '-outline'}` : 'md-map';
            break;
          case 'Chats':
            iconName =
              Platform.OS === 'ios' ? `ios-chatboxes${focused ? '' : '-outline'}` : 'md-chatboxes';
        }
        return (
          <Ionicons
            name={iconName}
            size={30}
            color={tintColor}
          />
        );
      },
      gesturesEnabled: false,
      headerRight: (
        <Button
          transparent
          onPress={() => navigation.navigate('SettingsScreen')}
          style={{ padding: 20, alignSelf: 'center' }}
        >
          <Ionicons
            name={Platform.OS === 'ios' ? 'ios-settings' : 'md-settings'}
            size={30}
            color='white'
          />
        </Button>
      ),
    }),
    tabBarOptions: {
      activeTintColor: Colors.tabIconSelected,
    },
    tabBarComponent: TabBarBottom,
    tabBarPosition: 'bottom',
    animationEnabled: false,
    swipeEnabled: false,
  }
);

class MainTabNavigation extends React.Component {
  static router = MainTabNavigator.router;

  componentWillUnmount() {
    this._notificationSubscription && this._notificationSubscription.remove();
  }

  render() {
    return (
      <MainTabNavigator
        navigation={this.props.navigation}
        screenProps={this.props.screenProps}
      />
    );
  }

  componentDidMount() {
    this._registerForPushNotifications();
    this.props.registerFetchThreadListAsync();
  }

  async _registerForPushNotifications() {
    if (!Constants.isDevice) {
      return;
    }

    // Get user permissions
    let { status } = await Permissions.askAsync(Permissions.NOTIFICATIONS);

    // Stop here if the user did not grant permissions
    if (status !== 'granted') {
      return;
    }

    // Get the token that uniquely identifies this device
    // TODO check if error when logged out
    let token = await Notifications.getExpoPushTokenAsync();
    modifyFirebaseUserState('expoPushToken', token);

    // Watch for incoming notifications
    this._notificationSubscription = Notifications.addListener(this._handleNotification);
  }

  _handleNotification = ({ origin, data }) => {
    console.log(`Push notification ${origin} with data: ${JSON.stringify(data)}`);
  };
}

const mapDispatchToProps = {
  registerFetchThreadListAsync
};

export default connect(
  null,
  mapDispatchToProps
)(MainTabNavigation);
