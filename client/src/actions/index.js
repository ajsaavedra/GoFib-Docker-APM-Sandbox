import axios from 'axios';
import {
    GET_FIB_VALUES,
    DELETE_FIB,
    ERROR
} from './types';

const ROOT_URL = 'http://localhost:80/api';

export function getFibValues() {
    return function(dispatch) {
        axios.get(`${ROOT_URL}/all`)
            .then(res => dispatch({ type: GET_FIB_VALUES, payload: res.data }))
            .catch(() => dispatch({ type: ERROR, payload: 'Something went wrong with your request' }));
    }
};

export function deleteFib(fib) {
    return function(dispatch) {
        axios.delete(`${ROOT_URL}/${fib}`)
            .then(res => dispatch({ type: DELETE_FIB, payload: res.data }))
            .catch(() => dispatch({ type: ERROR, payload: 'Could not delete fib value' }));
    }
};

export function putFibValue(value) {
    return function (dispatch) {
        axios.post(`${ROOT_URL}/fib`, { value: parseInt(value) })
            .then(() => console.log(`finished fetching {value}`))
            .catch(() => dispatch({ type: ERROR }))
    }
};
