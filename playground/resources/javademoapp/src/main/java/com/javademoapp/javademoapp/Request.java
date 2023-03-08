package com.javademoapp.javademoapp;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.ArrayList;

public class Request {

    @JsonProperty("request")
    private ArrayList<SubrequestChain> chains;

    public Request() {
        this.chains = new ArrayList<>();
    }
    public Request(ArrayList<SubrequestChain> chains) {
        this.chains = chains;
    }

    public ArrayList<SubrequestChain> getChains() {
        return chains;
    }
    public void setChains(ArrayList<SubrequestChain> chains) {
        this.chains = chains;
    }
    public void addChain(SubrequestChain chain) {
        this.chains.add(chain);
    }
}
