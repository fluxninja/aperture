package com.javademoapp.javademoapp;

import com.fasterxml.jackson.annotation.JsonProperty;

public class Subrequest {

    private String key;
    @JsonProperty("destination")
    private String destination;
    public Subrequest(String key, String destination) {
        this.key = key;
        this.destination = destination;
    }
    public Subrequest(String destination) {
        this.destination = destination;
    }

    public String getDestination() {
        return destination;
    }

    public void setDestination(String destination) {
        this.destination = destination;
    }

}
