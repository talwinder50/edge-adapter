#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@all
@rp_adapter
Feature: RP Adapter
  Background: Setup External Agent
    Given the "Mock Wallet" is running on "localhost" port "9081" with controller "http://localhost:9082"
    And the "Mock Issuer Adapter" is running on "localhost" port "10010" with controller "http://localhost:10011"

  Scenario: Register relying party
    When a request is sent to create an RP tenant with label "test-tenant"
    Then the trustbloc DID of the tenant with label "test-tenant" is resolvable
    And the client ID of the tenant with label "test-tenant" is registered at hydra

  Scenario: did-exchange with the web wallet
    Given a registered rp tenant with label "didexchange"
    When the rp tenant "didexchange" redirects the user to the rp adapter with scope "CreditCardStatement"
    And the rp adapter "didexchange" submits a CHAPI request to "Mock Wallet" with presentation-definitions and a didcomm invitation to connect
    And "Mock Wallet" accepts the didcomm invitation
    Then "Mock Wallet" connects with the RP adapter "didexchange"

  Scenario: Returns data from the user's issuer adapter to the relying party with the user's consent
    Given "Mock Issuer Adapter" and "Mock Wallet" have a didcomm connection
    And an rp tenant with label "userdata" that requests the "CreditCardStatement" scope from the "Mock Wallet"
    And the "Mock Wallet" provides a consent credential via CHAPI that contains the DIDs of rp "userdata" and issuer "Mock Issuer Adapter"
    When "Mock Issuer Adapter" responds to "userdata" with the user's data
    Then the user is redirected to the rp tenant "userdata"
    And the rp tenant "userdata" retrieves the user data from the rp adapter
