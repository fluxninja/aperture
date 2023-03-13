package com.javademoapp.javademoapp;

import java.util.ArrayList;
import java.util.List;

public class SubrequestChain {

    private List<Subrequest> subrequest;

    public SubrequestChain() {
        this.subrequest = new ArrayList<>();
    }

    public SubrequestChain(List<Subrequest> subrequest) {
        this.subrequest = subrequest;
    }

    public List<Subrequest> getSubrequest() {
        return subrequest;
    }

    public void setSubrequest(List<Subrequest> subrequest) {
        this.subrequest = subrequest;
    }

    public void addSubrequest(Subrequest subrequest) {
        this.subrequest.add(subrequest);
    }
}
