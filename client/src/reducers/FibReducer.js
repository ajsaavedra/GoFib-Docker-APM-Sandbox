import {
    GET_FIB_VALUES,
    ERROR
} from '../actions/types';

export default function(state = {}, action) {
    switch(action.type) {
        case GET_FIB_VALUES:
            return { ...state, values: action.payload, error: false };
        case ERROR:
            return { ...state, error: true }
    }
    return state;
};
