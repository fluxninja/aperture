package com.javademoapp.javademoapp;

import java.util.ArrayList;
import java.util.List;

public class Request {

    private List<List<Subrequest>> request;

    public Request() {
        this.request = new ArrayList<List<Subrequest>>();
    }

    public Request(List<List<Subrequest>> request) {
        this.request = request;
    }

    public List<List<Subrequest>> getRequest() {
        return request;
    }

    public void setRequest(List<List<Subrequest>> request) {
        this.request = request;
    }

    public void addRequest(List<Subrequest> chain) {
        this.request.add(chain);
    }
}
