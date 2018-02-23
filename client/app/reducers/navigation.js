import { NavigationActions } from 'react-navigation';

import { AppNavigator } from '../navigation';

import {
  CREATE_PROFILE,
  LOGIN,
  LOGOUT
} from '../actions/actionTypes';

const actionForLoggedOut = AppNavigator.router.getActionForPathAndParams(
  'Login'
);

const stateForLoggedOut = AppNavigator.router.getStateForAction(
  actionForLoggedOut
);

const actionForLoggedIn = AppNavigator.router.getActionForPathAndParams(
  'Main/List'
);

const stateForLoggedIn = AppNavigator.router.getStateForAction(
  actionForLoggedIn,
  stateForLoggedOut
);

const initialState = { stateForLoggedOut, stateForLoggedIn };

export default navigation = (state = initialState, action) => {

  switch (action.type) {

    case CREATE_PROFILE:
      return {
        ...state,
        stateForLoggedOut: AppNavigator.router.getStateForAction(
          AppNavigator.router.getActionForPathAndParams('CreateProfile'),
          stateForLoggedOut
        )
      };

    case LOGIN:
      return {
        ...state,
        stateForLoggedIn: AppNavigator.router.getStateForAction(
          actionForLoggedIn,
          state.stateForLoggedOut
        )
      };

    case LOGOUT:
      return {
        ...state,
        stateForLoggedOut: AppNavigator.router.getStateForAction(
          NavigationActions.reset({
            index: 0,
            actions: [NavigationActions.navigate({ routeName: "Login" })]
          })
        )
      };

    // case "Navigation/BACK":
    //   return {
    //     ...state,
    //     stateForLoggedOut: AppNavigator.router.getStateForAction(
    //       NavigationActions.back(),
    //       stateForLoggedOut
    //     )
    //   };

    default:
      return {
        ...state,
        stateForLoggedIn: AppNavigator.router.getStateForAction(
          action,
          state.stateForLoggedIn
        )
      };
  }
};
