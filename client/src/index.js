import React from 'react';
import { render } from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import reduxPromise from 'redux-promise';
import reduxThunk from 'redux-thunk';
import reducers from './reducers';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';

const createStoreWithMiddleware = applyMiddleware(reduxPromise, reduxThunk)(createStore);
const store = createStoreWithMiddleware(reducers);

render(
    <Provider store={store}>
        <BrowserRouter>
            <App />
        </BrowserRouter>
    </Provider>,
    document.getElementById('root')
);

serviceWorker.unregister();
