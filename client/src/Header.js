import React, { Component } from 'react';
import { Navbar, Nav, Form, FormControl, Button } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import logo from './logo.png';

class Header extends Component {
   render() {
        return (
            <Navbar bg="light" expand="lg">
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
                <Navbar.Brand>
                    <Link to="/">
                    <img
                        src={logo}
                        height="60"
                        className="d-inline-block align-top"
                    />
                    </Link>
                </Navbar.Brand>
                <Nav className="mr-auto">
                    <Link to="/all-values">All Values</Link>
                </Nav>
                <Form inline>
                    <FormControl type="text" placeholder="fib(n)" className="mr-sm-2" />
                    <Button variant="outline-dark">Calculate</Button>
                </Form>
            </Navbar.Collapse>
            </Navbar>
        );
    }
}

export default Header;
