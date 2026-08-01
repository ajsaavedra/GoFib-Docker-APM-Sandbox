import React, { Component } from 'react';
import { Form, FormControl, Button } from 'react-bootstrap';
import { connect } from 'react-redux';
import TableValues from './TableValues';
import { putFibValue } from '../actions';

class HomePage extends Component {
    constructor(props) {
        super(props);
        this.state = { fib: '', alert: null };

        this.handleChange = this.handleChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.renderAlert = this.renderAlert.bind(this);
    }

    handleChange(event) {
        this.setState({ fib: event.target.value });
    }

    onSubmit(event) {
        event.preventDefault();
        const fib = this.state.fib;
        if (fib === '' || isNaN(fib) || fib < 0) {
            return this.setState({
                fib: '',
                alert: 'Please enter a valid number >= 0'
            });
        }
        this.setState({ alert: null, fib: '' });
        this.props.putFibValue(fib);
    }

    renderAlert() {
        if (this.state.alert) {
            return (
                <div className="err">{this.state.alert}</div>
            );
        }
        return (null);
    }

    render() {
        return (
            <div className="main-content">
                <h1>Fibonacci Calculator</h1>
                <div>
                    <Form inline onSubmit={this.onSubmit}>
                        <div>
                            <FormControl value={this.state.fib} type="text" placeholder="fib(n)" className="mr-sm-2" onChange={this.handleChange} />
                            <Button variant="outline-dark" type="submit">Calculate</Button>
                        </div>
                        {this.renderAlert()}
                    </Form>
                    <TableValues/>
                </div>
            </div>
        );
    }
}

export default connect(null, { putFibValue })(HomePage);
