import {
  createReactNavigationReduxMiddleware,
  createReduxBoundAddListener,
} from 'react-navigation-redux-helpers';
import thunk from 'redux-thunk';

const middleware = [
  thunk,
  createReactNavigationReduxMiddleware(
    'root',
    state => state.nav,
  )
];
const addListener = createReduxBoundAddListener('root');

export {
  middleware,
  addListener,
};
