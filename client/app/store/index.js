import { createStore, applyMiddleware } from 'redux';

import { middleware } from './middleware';
import reducers from '../reducers';

export default createStore(reducers, applyMiddleware(...middleware));
