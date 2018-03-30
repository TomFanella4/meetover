import Expo, { Permissions } from 'expo';
import { times } from 'lodash';
import { serverURI } from '../constants/Common';

import {
  FETCH_MATCHES
} from './actionTypes';

const useMocks = false;

const fetchMatches = matches => ({
  type: FETCH_MATCHES,
  matches
});

export const fetchMatchesAsync = userId => {
  return async dispatch => {
    let matches;

    if (useMocks) {
      const uri = `${serverURI}/test/profile`;
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      let profile = await response.json();
      profile = JSON.parse(profile); // TODO make all the JSON.parse unnecessary

      matches = times(10, () => Object.assign({}, profile));
    } else {
      const { status } = await Permissions.askAsync(Permissions.LOCATION);
      if (status === 'granted') {
        const location = await Expo.Location.getCurrentPositionAsync({});
        const uri = `${serverURI}/match/${userId}`;
        const init = {
          method: 'POST',
          body: JSON.stringify({
            lat: location.coords.latitude,
            long: location.coords.longitude,
            timestamp: location.timestamp
          })
        };

        const response = await fetch(uri, init)
          .catch(err => console.log(err));
        const result = await response.json()
          .catch(err => console.log(err));
        matches = result.matches;

        matches.sort((a, b) => {
          if (a.distance > b.distance) return 1;
          if (a.distance < b.distance) return -1;
          return 0;
        });
        matches = matches.map(match => ({
          ...match.profile,
          location: {
            latitude: match.location.lat,
            longitude: match.location.long
          }
        }));
      }
    }

    dispatch(fetchMatches(matches));
  };
};
