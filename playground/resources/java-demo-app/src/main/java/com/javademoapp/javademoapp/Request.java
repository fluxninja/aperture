package com.javademoapp.javademoapp;

import java.util.ArrayList;
import java.util.List;

public class Request {

    private List<SubrequestChain> request;

    public Request() {
        this.request = new ArrayList<>();
    }

    public Request(List<SubrequestChain> request) {
        this.request = request;
    }

    public List<SubrequestChain> getRequest() {
        return request;
    }

    public void setRequest(List<SubrequestChain> request) {
        this.request = request;
    }

    public void addRequest(SubrequestChain chain) {
        this.request.add(chain);
    }
}
