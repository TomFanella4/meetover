import React from 'react';
import { Platform } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { TabNavigator, TabBarBottom } from 'react-navigation';
import { Button } from 'native-base';
import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { registerFetchThreadListAsync } from '../actions/chatActions';

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

  render() {
    return (
      <MainTabNavigator
        navigation={this.props.navigation}
        screenProps={this.props.screenProps}
      />
    );
  }

  componentDidMount() {
    this.props.registerFetchThreadListAsync();
  }
}

const mapDispatchToProps = {
  registerFetchThreadListAsync
};

export default connect(
  null,
  mapDispatchToProps
)(MainTabNavigation);
