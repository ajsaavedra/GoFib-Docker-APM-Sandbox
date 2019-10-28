import React, { Component } from 'react';
import { Navbar } from 'react-bootstrap';
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
                        alt="arf!"
                    />
                    </Link>
                </Navbar.Brand>
            </Navbar.Collapse>
            </Navbar>
        );
    }
}

export default Header;
