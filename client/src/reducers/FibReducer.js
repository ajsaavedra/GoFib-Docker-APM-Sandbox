import {
    GET_FIB_VALUES,
    ERROR,
    DELETE_FIB
} from '../actions/types';

export default function(state = {}, action) {
    switch(action.type) {
        case GET_FIB_VALUES:
            return { ...state, values: action.payload, error: false };
        case DELETE_FIB:
            return { ...state, deleted: true, value: action.payload };
        case ERROR:
            return { ...state, error: true };
    }
    return state;
};
