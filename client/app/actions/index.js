import {
  FETCH_MATCHES,
  LOGIN,
  LOGOUT,
  UPDATE_USER_PROFILE,
} from './actionTypes';

export const fetchMatches = matches => ({
  type: FETCH_MATCHES,
  matches
});

export const fetchMatchesAsync = userId => {
  return async dispatch => {
    const uri = `https://meetover.herokuapp.com/match/${userId}`;
    const init = { method: 'POST' };

    // const response = await fetch(uri, init);
    // const matches = await response.json();

    const matches = ['Matt', 'Tom', 'Colin', 'Ryan', 'Krutarth', 'Austin', 'Alec', 'Janet', 'Doug', 'Joey', 'Mark'];

    dispatch(fetchMatches(matches));
  };
};

export const login = userProfile => ({
  type: LOGIN,
  userProfile
});

export const logout = () => ({
  type: LOGOUT
});

export const modifyUserProfile = userProfile => ({
  type: UPDATE_USER_PROFILE,
  userProfile
});
