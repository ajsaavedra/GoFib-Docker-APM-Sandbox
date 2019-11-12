import { combineReducers } from 'redux';
import FibReducer from './FibReducer';

const rootReducer = combineReducers({
    fib: FibReducer,
});

export default rootReducer;
