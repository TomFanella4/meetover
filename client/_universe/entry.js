import Expo from 'expo';
import App from '../app/App';

if (__DEV__) {
  Expo.KeepAwake.activate();
}

Expo.registerRootComponent(App);
