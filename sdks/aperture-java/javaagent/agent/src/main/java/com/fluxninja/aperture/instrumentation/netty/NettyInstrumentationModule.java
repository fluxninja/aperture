package com.fluxninja.aperture.instrumentation.netty;

import com.fluxninja.aperture.instrumentation.InstrumentationModule;
import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;
import java.util.ArrayList;
import java.util.List;

public class NettyInstrumentationModule implements InstrumentationModule {
    @Override
    public List<TransformerInstrumentation> getTransformers() {
        return new ArrayList<TransformerInstrumentation>() {
            {
                add(new NettyServerInstrumentation());
            }
        };
    }
}
