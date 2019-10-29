import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import FibReducer from './FibReducer';

const rootReducer = combineReducers({
    fib: FibReducer,
    form: formReducer
});

export default rootReducer;
