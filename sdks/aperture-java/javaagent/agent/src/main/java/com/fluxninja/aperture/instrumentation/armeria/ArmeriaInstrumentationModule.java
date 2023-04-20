package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.instrumentation.InstrumentationModule;
import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;
import java.util.ArrayList;
import java.util.List;

public class ArmeriaInstrumentationModule implements InstrumentationModule {
    @Override
    public List<TransformerInstrumentation> getTransformers() {
        return new ArrayList<TransformerInstrumentation>() {
            {
                add(new ArmeriaClientInstrumentation());
                add(new ArmeriaServerInstrumentation());
            }
        };
    }
}
