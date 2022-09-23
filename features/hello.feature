Feature: greet a user
  I need to be able to receive a greet with my name

  Scenario: Greet User
    Given Hello service is running
    When call Hello with username "Jhon"
    Then Hello should return message "Hello Jhon"